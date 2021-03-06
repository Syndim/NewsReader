(function() {
  var commentTemplate, displayComments, getComment;

  commentTemplate = '<div class="well well-sm"> {{{parentComment}}} <p>{{commentContent}}</p> <p><small>{{commentInfo}}</small></p> </div>';

  getComment = function(template, cmt, cmtList, cmtStore) {
    var comment, currentComment, innerComment, parentComment;
    if (!cmt || !cmtStore[cmt.tid]) {
      return "";
    }
    parentComment = [
      (function() {
        var _i, _len, _results;
        _results = [];
        for (_i = 0, _len = cmtList.length; _i < _len; _i++) {
          comment = cmtList[_i];
          if (cmt.parent && comment.tid === cmt.parent) {
            _results.push(comment);
          }
        }
        return _results;
      })()
    ][0][0];
    innerComment = getComment(template, parentComment, cmtList, cmtStore);
    currentComment = cmtStore[cmt.tid];
    return template({
      'parentComment': innerComment,
      'commentContent': $('<div/>').html(currentComment.comment).text(),
      'commentInfo': currentComment.host_name + " " + currentComment.name + " " + currentComment.date + " 支持(" + currentComment.score + ") 反对(" + currentComment.reason + ")"
    });
  };

  displayComments = function(comments) {
    var cmtResult, comment, hotCmtResult, template, _i, _len;
    for (_i = 0, _len = comments.length; _i < _len; _i++) {
      comment = comments[_i];
      if (comment) {
        comments = JSON.parse(comment);
      }
    }
    template = Handlebars.compile(commentTemplate);
    cmtResult = [
      (function() {
        var _j, _len1, _ref, _results;
        _ref = comments.cmntlist;
        _results = [];
        for (_j = 0, _len1 = _ref.length; _j < _len1; _j++) {
          comment = _ref[_j];
          _results.push(getComment(template, comment, comments.cmntlist, comments.cmntstore));
        }
        return _results;
      })()
    ][0];
    hotCmtResult = [
      (function() {
        var _j, _len1, _ref, _results;
        _ref = comments.hotlist;
        _results = [];
        for (_j = 0, _len1 = _ref.length; _j < _len1; _j++) {
          comment = _ref[_j];
          _results.push(getComment(template, comment, comments.hotlist, comments.cmntstore));
        }
        return _results;
      })()
    ][0];
    $('#blog-comment').html(cmtResult.join("<hr />"));
    $('#blog-hotlist').html(hotCmtResult.join(''));
    return console.log(comments);
  };

  $(document).ready(function() {
    return displayComments(comments);
  });

}).call(this);
