from django.contrib import admin
from .models import *

# Register your models here.

class NewsAdmin(admin.ModelAdmin):
    pass

class CommentAdmin(admin.ModelAdmin):
    pass

admin.site.register(News, NewsAdmin)
admin.site.register(Comment, CommentAdmin)
