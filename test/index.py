import os


def handler(*_):
    return sorted(os.environ)
