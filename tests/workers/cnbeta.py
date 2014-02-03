#coding: utf8

import unittest
import os
from workers import cnbeta

class CnBetaTests(unittest.TestCase):

    def setUp(self):
        dir_name = os.path.dirname(os.path.realpath(__file__))
        with open(os.path.join(dir_name, "testdata/cnbeta/1.txt"), encoding='utf8') as f:
            self.text_data = f.read()

    def test_get_comments_details(self):
        details = cnbeta.get_comments_details(self.text_data)
        self.assertIsNotNone(details)
        self.assertEqual(2, len(details))
        self.assertEqual("270774", details[0])
        self.assertEqual("4f299", details[1])

    def test_get_op_code(self):
        op_code = cnbeta.get_op_code("270774", "4f299", 1)
        self.assertEqual("MSwyNzA3NzQsNGYyOTk=1234567", op_code)

    def test_get_comments(self):
        op_code = "MSwyNzA2NTAsYTYwMWI%3DuscCHONc"
        comments = cnbeta.get_comments(op_code)
        self.assertIsNotNone(comments)

    def test_get_article_list(self):
        article_list = cnbeta.get_article_list(1)
        self.assertIsNotNone(article_list)
        self.assertEqual(30, len(article_list))
