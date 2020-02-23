from .base import *  # noqa

DEBUG = True
ALLOWED_HOSTS = ["localhost", "127.0.0.1", "0.0.0.0"]
REST_FRAMEWORK["DEFAULT_PERMISSION_CLASSES"] = ("rest_framework.permissions.AllowAny",)
