from sentence_transformers import SentenceTransformer

# Load a pre-trained model
# all-MiniLM-L6-v2: A smaller, faster model that is suitable for many tasks and works well with limited resources.
# all-distilroberta-v1: A distilled version of RoBERTa optimized for performance and efficiency.
# paraphrase-MiniLM-L6-v2: Optimized for generating paraphrase embeddings.
model_name = 'paraphrase-MiniLM-L6-v2'  # You can choose other models from Sentence Transformers
model = SentenceTransformer(model_name)

# Define your text
texts = ["This is a sample sentence.", "Here is another sentence."]

# Generate embeddings
embeddings = model.encode(texts)

# Print the embeddings
for i, emb in enumerate(embeddings):
    print(f"Embedding for text {i} [vec_len: {len(emb)}]: {emb}")