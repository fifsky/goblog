//评论
(function($){

    var get_max_days = function(year, month) {
        switch (month) {
            case 1:
            case 3:
            case 5:
            case 7:
            case 8:
            case 10:
            case 12:
                return 31;
            case 4:
            case 6:
            case 9:
            case 11:
                return 30;
            case 2:
                if (year % 400 == 0 || (year % 4 == 0 && year % 100 != 0)) {
                    return 29;
                }
                else {
                    return 28;
                }
            default:
                return 0;
        }
    };

    var get_week = function(y, m, d) {
        return (new Date(y+'/'+m+'/'+d)).getDay();
    };


    function calendar(y,m) {
        var days = get_max_days(y, m);
        var start_week = get_week(y, m, 1);
        var html = '<table class="canlender"><tr><td><a href="javascript:void(0);" class="_cal_handel _prev_year">&lt;&lt;</a></td><td><a href="javascript:void(0);" class="_cal_handel _prev_month">&lt;</a></td><td colspan="3">'+y+'-'+(m > 9 ? m : '0'+m) +'</td><td><a href="javascript:void(0);" class="_cal_handel _next_month">&gt;</a></td><td><a href="javascript:void(0);" class="_cal_handel _next_year">&gt;&gt;</a></td></tr><tr><td class="happy">日</td><td>一</td><td>二</td><td>三</td><td>四</td><td>五</td><td class="happy">六</td></tr>';
        for (var i = 0,d = 1; i <= 41; i++) {
            if(i%7 == 0){
                html += '<tr>';
            }

            if(i < start_week || d > days){
                html += '<td></td>';
            }else if(d <= days){
                var today = new Date();
                var cy = today.getFullYear(),cm = today.getMonth() + 1,cd = today.getDate();

                if(y === cy && cm === m && d === cd){
                    html += '<td class="today">'+(d++)+'</td>';
                }else{
                    html += '<td>'+(d++)+'</td>';
                }
            }

            if((i+1)%7 == 0){
                html += '</tr>';
                if(d > days){
                    break;
                }
            }
        }
        html += "</table>";

        var $cal = $(html);

        $cal.find('a._cal_handel').bind('click',function(){
            var year = y,month = m;
            if($(this).hasClass('_prev_year')){
                year--;
            }else if($(this).hasClass('_prev_month')){
                month--;
                if(month < 1){
                    month = 12;
                    year--;
                }
            }else if($(this).hasClass('_next_month')){
                month++;
                if(month > 12){
                    month = 1;
                    year++;
                }
            }else if($(this).hasClass('_next_year')){
                year++;
            }

            calendar(year,month);
            return false;
        });

        $("#blog-calendar").html($cal);
    }

    if($('#blog-calendar').length){
        var today = new Date(),year = today.getFullYear(),month = today.getMonth() + 1;
        calendar(year,month);
    }

    $('#comment_form').submit(function(){

        var $el = $(this).find('input[name="name"]');
        if($.trim($el.val()) == ''){
            $el.focus().closest('p').addClass('error');
            return false;
        }

        // $el = $(this).find('input[name="comment_email"]');
        // if(($el.data('require') && $.trim($el.val()) == '') || !/^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/.test($el.val())){
        //     $el.focus().closest('p').addClass('error');
        //     return false;
        // }

        // $el = $(this).find('input[name="comment_url"]');
        // if(($el.data('require') && $.trim($el.val()) == '') || !/^https?:\/\/[A-Za-z0-9]+\.[A-Za-z0-9]+[\/=\?%\-&_~`@[\]\':+!]*([^<>\"\"])*$/.test($el.val())){
        //     $el.focus().closest('p').addClass('error');
        //     return false;
        // }

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
                    var $content = $(data.data.content).hide();
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
})(jQuery);