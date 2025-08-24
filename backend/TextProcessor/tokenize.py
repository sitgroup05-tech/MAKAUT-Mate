from transformers import DPRContextEncoder, DPRContextEncoderTokenizer
import torch

import numpy as np
import random
from transformers import DPRQuestionEncoder, DPRQuestionEncoderTokenizer
from transformers import AutoTokenizer, AutoModelForCausalLM

# from collections import fil


context_tokenizer = DPRContextEncoderTokenizer.from_pretrained('facebook/dpr-ctx_encoder-single-nq-base')
context_encoder = DPRContextEncoder.from_pretrained('facebook/dpr-ctx_encoder-single-nq-base')

def process_text(context : str):
    
    # divding into chunck
    paras = context.split('\n')

    paras = list(filter(lambda x : len(x) > 0 , paras))

    def encode_contexts(text_list):
        # Encode a list of texts into embeddings
        embeddings = []
        for text in text_list:
            inputs = context_tokenizer(text, return_tensors='pt', padding=True, truncation=True, max_length=256)
            outputs = context_encoder(**inputs)
            embeddings.append(outputs.pooler_output)
        return torch.cat(embeddings).detach().numpy()

    context_embeddings = encode_contexts(paras)
    print(context_embeddings.shape)
    return context_embeddings