from rest_framework import viewsets

from dumbwaiter import models as m, serializers as s


class DayViewSet(viewsets.ModelViewSet):
    queryset = m.Day.objects.all()
    serializer_class = s.DaySerializer

