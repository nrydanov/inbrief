from pydantic_settings import BaseSettings, SettingsConfigDict

class ModelSettings(BaseSettings):
    model_name: str = "sentence-transformers/all-MiniLM-L6-v2"
    language: str = 'russian'
    trust_remote_code: bool = True

class RedisSettings(BaseSettings):
    host: str = "127.0.0.1"
    port: int = 6379

class MinioSettings(BaseSettings):
    endpoint: str = "127.0.0.1:9000"
    access_key: str = "minioadmin"
    secret_key: str = "minioadmin"
    bucket: str = "inbrief"

class AppSettings(BaseSettings):
    redis: RedisSettings = RedisSettings()
    minio: MinioSettings = MinioSettings()
    embedding_model: ModelSettings = ModelSettings()
    port: int = 8081

    model_config = SettingsConfigDict(
        env_prefix="APP_",
        extra="ignore",
        env_nested_delimiter="_",
        env_file='.env',
    )


settings = AppSettings()
