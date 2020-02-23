from django.conf import settings
from django.db import models


class Day(models.Model):

    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE)
    date = models.DateField()
    am_weight = models.FloatField(default=0)
    pm_weight = models.FloatField(default=0)
    snack = models.IntegerField(default=0)
    breakfast = models.IntegerField(default=0)
    lunch = models.IntegerField(default=0)
    dinner = models.IntegerField(default=0)
    exercise = models.IntegerField(default=0)

    class Meta:
        unique_together = ("user", "date")
