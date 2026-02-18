

import os
from dotenv import load_dotenv , set_key , find_dotenv

load_dotenv()

env_path = find_dotenv()


DB_URL : str = os.getenv("DB_URL")


if __name__ == "__main__" :
    set_key(env_path,"PYTHON_RPC_PID",str(os.getpid()))
    set_key(env_path,"PYTHON_RPC_ACTIVE",str(True))
    print(os.getenv("PGHOST"))
    print(os.getenv("PYTHON_RPC_ACTIVE"))