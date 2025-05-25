import logging

import numpy as np
import redis
from minio import Minio
from umap import UMAP
from config.settings import settings
from sentence_transformers import SentenceTransformer
from bertopic import BERTopic
from bertopic.representation import KeyBERTInspired
import json

logger = logging.getLogger(__name__)

def clustering(queue):
    logger.setLevel(logging.DEBUG)

    logger.info("Starting clustering process")
    r = redis.Redis(host=settings.redis.host, port=settings.redis.port, db=0)
    if not r.ping():
        logger.error("Failed to connect to redis")
        return
    else:
        logger.info("Connected to redis")

    s3 = Minio(
        settings.minio.endpoint,
        access_key=settings.minio.access_key,
        secret_key=settings.minio.secret_key,
        secure=False,
    )

    model = SentenceTransformer(
        settings.embedding_model.model_name,
        trust_remote_code=settings.embedding_model.trust_remote_code
    )

    bertopic = BERTopic(
        language=settings.embedding_model.language,
        embedding_model=model,
        representation_model=KeyBERTInspired()
    )


    cached_entities = r.get("clustering:entities")
    if cached_entities is not None:
        entities = json.loads(cached_entities)
    else:
        entities = []


    pubsub = r.pubsub()
    pubsub.subscribe("inbrief")

    dim = model.encode("Hello world!", task="separation").shape[0]
    logger.debug("Embedding dimension: %s", dim)

    umap = UMAP(n_neighbors=15, n_components=2, metric='cosine')

    for msg in pubsub.listen():
        try:
            try:
                if queue.get(block=False) == None:
                    logger.debug("Got stop signal")
                    break
            except:
                pass
            if msg is None:
                continue

            logger.debug("Received message: %s", msg)
            if msg["type"] != "message":
                continue

            filename = f"{msg['data'].decode('utf-8')}.json"
            resp = s3.get_object("inbrief", filename)
            payload = resp.json()

            new_texts = list(map(lambda x: x['text'], payload))
            logger.debug("Encoding texts")
            new_embeddings = model.encode(new_texts, task="separation").tolist()

            for i, entity in enumerate(payload):
                entity['embedding'] = new_embeddings[i]

            entities.extend(payload)

            texts = list((map(lambda x: x['text'], entities)))
            embeddings = np.array(list(map(lambda x: x['embedding'], entities)))

            try:
                logger.debug("Reducing embeddings")
                reduced_embeddings = umap.fit_transform(embeddings)
                logger.debug("Fitting bertopic")
                bertopic.fit_transform(texts, embeddings=embeddings)

                logger.debug("Writing topics")
                bertopic.visualize_topics().write_html("static/topics.html")
                logger.debug("Writing documents")
                bertopic.visualize_documents(texts, reduced_embeddings=reduced_embeddings).write_html("static/documents.html")
            except Exception as e:
                logger.warning(f"Got an error while clustering, will try later: {e}", exc_info=True)
                continue

            logger.debug("Number of texts: %s", len(texts))
            dumped_entities = json.dumps(entities)
            r.set("clustering:entities", dumped_entities)
            dumped_topics = json.dumps(bertopic.topics_)
            r.set("clustering:topics", dumped_topics)
        # NOTE(nrydanov): Change my mind
        except Exception as e:
            logger.error(f"Error in clustering service: {e}", exc_info=True)
