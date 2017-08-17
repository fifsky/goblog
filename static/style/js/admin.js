(function($){

    $('a[rel="confirmTodo"]').click(function(){
        var title = $(this).attr('title') ? $(this).attr('title') : '确认？';
        return confirm(title);
    });

})(jQuery);