var baseApp = {
    init: function () {
        this.initAside()
        this.confirmDelete()
        this.resizeIframe()
    },
    initAside: function () {
        $('.aside h4').click(function () {
            $(this).siblings('ul').slideToggle();
        })
    },
    //设置iframe的高度
    resizeIframe: function () {
        $("#rightMain").height($(window).height() - 80)
    },
    // 删除提示
    confirmDelete: function () {
        $(".delete").click(function () {
            var flag = confirm("您确定要删除吗?")
            return flag
        })
    }
};

$(function () {
    baseApp.init();
});
