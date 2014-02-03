from django.conf.urls import patterns, include, url

from cnbeta.feed import LatestEntriesFeed

from django.contrib import admin
admin.autodiscover()

urlpatterns = patterns('',
    # Examples:
    url(r'^$', 'cnbeta.views.index', name='index'),
    url(r'^page/(?P<page>\d+)$', 'cnbeta.views.index', name='index'),
    url(r'^news/(?P<news_id>\d+)$', 'cnbeta.views.detail', name='detail'),
    url(r'^news/feed$', LatestEntriesFeed()),
    # url(r'^blog/', include('blog.urls')),

    url(r'^admin/', include(admin.site.urls)),
)
