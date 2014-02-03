from django.shortcuts import render
from django.core.paginator import Paginator
from .models import *

# Create your views here.

def index(request, page=1):
    paged = Paginator(News.objects.all().order_by("-origin_id"), 10)
    return render(request, 'list.html', {
        'paged': paged.page(page)
    })

def detail(request, news_id):
    news = News.objects.get(pk=news_id)
    comment = Comment.objects.get(
            origin_site=news.origin_site,
            origin_id=news.origin_id
            )
    return render(request, 'detail.html', {
        'news': news,
        'comment': comment
        })
