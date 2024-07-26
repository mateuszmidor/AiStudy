from sentence_transformers import SentenceTransformer
import sys
import json

# Load a pre-trained model
# all-MiniLM-L6-v2: A smaller, faster model that is suitable for many tasks and works well with limited resources.
# all-distilroberta-v1: A distilled version of RoBERTa optimized for performance and efficiency.
# paraphrase-MiniLM-L6-v2: Optimized for generating paraphrase embeddings.
model_name = 'paraphrase-MiniLM-L6-v2'  # 384 dimensions
model = SentenceTransformer(model_name)

# Usage: python main.py "My text to embed"
def main():
    text = sys.argv[1]
    sentences = [text]

    embeddings = model.encode(sentences)
    embeddings_list = embeddings.tolist()[0]

    print(embeddings_list)

if __name__ == "__main__":
    main()