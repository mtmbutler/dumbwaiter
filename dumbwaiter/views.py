from rest_framework import viewsets

from dumbwaiter import models as m, serializers as s


class DayViewSet(viewsets.ModelViewSet):
    serializer_class = s.DaySerializer

    def get_queryset(self):
        if self.request.user.is_authenticated:
            return m.Day.objects.filter(user=self.request.user)
        return m.Day.objects.none()

    def perform_create(self, serializer):
        serializer.save(user=self.request.user)
