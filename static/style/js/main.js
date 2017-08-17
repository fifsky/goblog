//评论
(function($){
    $('#comment_form').submit(function(){

        var $el = $(this).find('input[name="comment_author"]');
        if($.trim($el.val()) == ''){
            $el.focus().closest('p').addClass('error');
            return false;
        }

        $el = $(this).find('input[name="comment_email"]');
        if(($el.data('require') && $.trim($el.val()) == '') || !/^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/.test($el.val())){
            $el.focus().closest('p').addClass('error');
            return false;
        }

        $el = $(this).find('input[name="comment_url"]');
        if(($el.data('require') && $.trim($el.val()) == '') || !/^https?:\/\/[A-Za-z0-9]+\.[A-Za-z0-9]+[\/=\?%\-&_~`@[\]\':+!]*([^<>\"\"])*$/.test($el.val())){
            $el.focus().closest('p').addClass('error');
            return false;
        }

        var self = this;

        $.ajax({
            'url': $(this).attr('action'),
            'type': this.method,
            'data': $(this).serializeArray(),
            'dataType': 'json',
            beforeSend: function(xhr) {
                $(self).find(':submit').addClass('formdisabled').attr('disabled',true);
            },
            complete:function(){
                $(self).find(':submit').removeClass('formdisabled').attr('disabled',false);
            },
            success: function(data) {
                if (data.statusCode == 200){
                    if($('#comment .comment-null').size()){
                        $('#comment .comment-null').remove();
                    }
                    var $content = $(data.content).hide();
                    $('#comment ul').append($content);
                    $content.fadeIn(200);
                    self.reset();
                }else{
                    alert(data.message);
                }
            }
        });
        return false;
    });
})(jQuery);
(function($){
    $(window).scroll(function(){
        if($(this).scrollTop() > 300){
            $('#scroll_top').show();
        }else{
            $('#scroll_top').hide();
        }
    });
    $('#scroll_top').click(function(){
        $.scrollTo($('body'),200);
    })
})(jQuery);