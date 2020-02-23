from django.conf.urls import url, include
from rest_framework import routers

from dumbwaiter.views import DayViewSet

router = routers.DefaultRouter()
router.register("days", DayViewSet)

urlpatterns = [
    url(r"^", include(router.urls)),
    url(r"^api-auth/", include("rest_framework.urls", namespace="rest_framework")),
]
