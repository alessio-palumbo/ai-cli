import requests
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

OLLAMA_URL = "http://localhost:11434/api/generate"
OLLAMA_MODEL = "llama3"


class GenerateRequest(BaseModel):
    prompt: str


@app.post("/generate")
def generate(req: GenerateRequest):
    payload = {"model": OLLAMA_MODEL, "prompt": req.prompt, "stream": False}
    return requests.post(OLLAMA_URL, json=payload).json()
