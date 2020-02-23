from django.contrib.auth.models import AnonymousUser
from rest_framework import viewsets

from dumbwaiter import models as m, serializers as s


class DayViewSet(viewsets.ModelViewSet):
    serializer_class = s.DaySerializer

    def get_queryset(self):
        user = self.request.user
        if isinstance(user, AnonymousUser):
            return m.Day.objects.none()
        return m.Day.objects.filter(user=user)

    def perform_create(self, serializer):
        serializer.save(user=self.request.user)
