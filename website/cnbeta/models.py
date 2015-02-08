from django.db import models

# Create your models here.

class News(models.Model):
    title = models.CharField(max_length=100)
    intro = models.CharField(max_length=3000)
    content = models.TextField(max_length=50000)
    created_at = models.CharField(max_length=1000)
    origin_id = models.IntegerField(db_index=True)
    origin_site = models.CharField(max_length=100)

class Comment(models.Model):
    content = models.TextField(max_length=50000)
    updated_at = models.CharField(max_length=1000)
    origin_id = models.IntegerField(db_index=True)
    origin_site = models.CharField(max_length=100)
