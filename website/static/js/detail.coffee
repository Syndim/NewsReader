commentTemplate = '
<div class="well well-sm">
    {{{parentComment}}}
    <p>{{commentContent}}</p>
    <p><small>{{commentInfo}}</small></p>
</div>
'

getComment = (template, cmt, cmtList, cmtStore) ->
    return "" if not cmt or not cmtStore[cmt.tid]
    parentComment = [comment for comment in cmtList when cmt.parent and comment.tid == cmt.parent][0][0]
    innerComment = getComment template, parentComment, cmtList, cmtStore
    currentComment = cmtStore[cmt.tid]
    return template {
        'parentComment': innerComment
        'commentContent': $('<div/>').html(currentComment.comment).text()
        'commentInfo': currentComment.host_name + " " + currentComment.name + " " + currentComment.date + " 支持(" + currentComment.score + ") 反对(" + currentComment.reason + ")"
    }

displayComments = (comments) ->
    comments = JSON.parse atob(comment).substr(6) for comment in comments when comment
    template = Handlebars.compile commentTemplate
    cmtResult = [getComment(template, comment, comments.cmntlist, comments.cmntstore) for comment in comments.cmntlist][0]
    hotCmtResult = [getComment(template, comment, comments.hotlist, comments.cmntstore) for comment in comments.hotlist][0]
    $('#blog-comment').html cmtResult.join("<hr />")
    $('#blog-hotlist').html hotCmtResult.join('')
    console.log comments

$(document).ready () ->
    displayComments comments

