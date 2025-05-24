import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
import multiprocessing
from clustering import clustering
import redis
from config.settings import settings
import json


logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s %(message)s"
)

# Глобальные переменные


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Управление жизненным циклом приложения"""

    queue = multiprocessing.Queue()
    process = multiprocessing.Process(target=clustering, args=(queue,))
    process.start()

    logger.info("FastAPI app started")

    app.state.r = redis.Redis(host=settings.redis.host, port=settings.redis.port, db=0)

    if not app.state.r.ping():
        logger.error("Failed to connect to redis")
        return
    else:
        logger.info("Connected to redis")
    yield

    queue.put(None)

    # Очистка при завершении
    logger.info("Shutting down...")
    process.join()

# Создаем FastAPI приложение с lifespan
app = FastAPI(lifespan=lifespan)

def get_entities(r):
    return json.loads(r.get("clustering:entities"))

def get_topics(r):
    return json.loads(r.get("clustering:topics"))

@app.get("/get_summary")
async def get_summary():
    r = app.state.r
    entities = get_entities(r)
    topics = get_topics(r)
    output = [
        {k: v for k, v in entity.items() if k != 'embeddings'}
        for entity in entities
    ]
    return {"entities": output, "labels": topics}

@app.get("/health")
async def health():
    return "fine"
