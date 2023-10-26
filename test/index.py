import json
import os


def handler(*_):
    print(json.dumps(sorted(os.environ)))
    return True
