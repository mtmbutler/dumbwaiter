from django.conf.urls import url, include
from rest_framework import routers

from dumbwaiter.views import DayViewSet
from server.views import UserViewSet

router = routers.DefaultRouter()
router.register("users", UserViewSet, basename="User")
router.register("days", DayViewSet, basename="Day")

urlpatterns = [
    url(r"^", include(router.urls)),
    url(r"^api-auth/", include("rest_framework.urls", namespace="rest_framework")),
]
