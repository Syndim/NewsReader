#coding:utf-8
from django.contrib.syndication.views import Feed
from django.utils.feedgenerator import Rss201rev2Feed

from .models import News

import datetime

class ExtendedRSSFeed(Rss201rev2Feed):
    mime_type = 'application/xml'
    """
    Create a type of RSS feed that has content:encoded elements.
    """
    def root_attributes(self):
        attrs = super(ExtendedRSSFeed, self).root_attributes()
        attrs['xmlns:content'] = 'http://purl.org/rss/1.0/modules/content/'
        return attrs

    def add_item_elements(self, handler, item):
        super().add_item_elements(handler, item)
        #handler.addQuickElement('content:encoded', item['content_encoded'])


class LatestEntriesFeed(Feed):
    feed_type = ExtendedRSSFeed

    # Elements for the top-level, channel.
    title = u"Newspaper"
    link = "http://np.syndim.org"
    author = 'Syndim'
    description = u"Purl news site, without ads"

    def items(self):
        return News.objects.all().order_by('-origin_id')[:10]

    # Elements for each item.
    def item_title(self, item):
        return item.title

    def item_description(self, item):
        return item.intro

    def item_author_name(self, item):
        return item.origin_site

    def item_pubdate(self, item):
        return datetime.datetime.strptime(item.created_at, "%Y-%m-%d %H:%M:%S")

    def item_content_encoded(self, item):
        return item.intro + item.content

    def item_link(self, item):
        return "http://np.syndim.org/news/" + str(item.id)
