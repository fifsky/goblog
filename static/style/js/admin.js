(function ($) {

    $('a[rel="confirmTodo"]').click(function () {
        var title = $(this).attr('title') ? $(this).attr('title') : '确认？';
        return confirm(title);
    });

    $('.all-selected').on("click", function () {
        $('.list').find('input[name="ids"]').prop("checked", true);
    });

    $('.inverse-selected').on("click", function () {
        var $el = $('.list').find('input[name="ids"]');
        var $el1 = $el.filter(":checked");
        var $el2 = $el.not(":checked");
        $el1.prop("checked", false);
        $el2.prop("checked", true);
    });

})(jQuery);