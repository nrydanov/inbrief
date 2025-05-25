import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
import multiprocessing
from clustering import clustering
import redis
from config.settings import settings
import json
from fastapi.middleware.cors import CORSMiddleware


logger = logging.getLogger(__name__)

@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.setLevel(logging.DEBUG)


    queue = multiprocessing.Queue()
    process = multiprocessing.Process(target=clustering, args=(queue,), daemon=True)
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
    logger.info("Shutting down...")

app = FastAPI(lifespan=lifespan)

# Добавляем CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173"],  # URL фронтенда
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

def get_entities(r):
    return json.loads(r.get("clustering:entities"))

def get_topics(r):
    return json.loads(r.get("clustering:topics"))

@app.get("/get_summary")
async def get_summary():
    r = app.state.r
    entities = get_entities(r)
    topics = get_topics(r)

    # Добавляем логирование для отладки
    logger.debug(f"Entities from Redis: {entities}")
    logger.debug(f"Topics from Redis: {topics}")

    # Создаем словарь для группировки сущностей по кластерам
    clusters = {}

    # Группируем сущности по их меткам кластеров
    for i, entity in enumerate(entities):
        # Получаем метку кластера для текущей сущности из массива topics
        cluster_label = topics[i] if i < len(topics) else None

        if cluster_label is not None:
            # Создаем копию сущности без embedding
            entity_without_embedding = {k: v for k, v in entity.items() if k != 'embedding'}

            # Добавляем сущность в соответствующий кластер
            if cluster_label not in clusters:
                clusters[cluster_label] = []
            clusters[cluster_label].append(entity_without_embedding)

    # Добавляем логирование для отладки
    logger.debug(f"Clusters after grouping: {clusters}")

    # Преобразуем словарь кластеров в список с данными в формате, ожидаемом фронтендом
    stories = []
    for i, (label, entities) in enumerate(clusters.items()):
        # Создаем timeline из сущностей, только если есть валидный ts
        timeline = [
            {
                "datetime": entity.get("ts"),
                "text": entity.get("text", "")
            }
            for entity in entities
            if entity.get("ts") not in (None, "")
        ]

        # Сортируем timeline по datetime (от новых к старым)
        timeline.sort(key=lambda x: x["datetime"], reverse=True)

        story = {
            "id": i + 1,
            "title": f"Кластер {i+1}",
            "summary": "Краткая сводка по всем новостям в этом сюжете...",
            "timeline": timeline,
            "tags": [str(label)] if label is not None else []
        }
        stories.append(story)

    # Добавляем логирование для отладки
    logger.debug(f"Final stories: {stories}")

    return {"stories": stories}

@app.get("/health")
async def health():
    return "fine"
