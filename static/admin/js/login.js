var app = {
    init: function() {
        this.getCaptcha();
        this.captchaImgChange();
    },
    getCaptcha: function() {
        $.get("/admin/captcha?t=" + Math.random(), function(response) {
            console.log(response);
            $("#captchaId").val(response.captchaId);
            $("#captchaImg").attr("src", response.captchaImg);
        }).fail(function(error) {
            console.error("Failed to get captcha:", error);
        });
    },
    captchaImgChange: function() {
        var that = this;
        $("#captchaImg").click(function() {
            that.getCaptcha();
        });
    }
};

$(function() {
    app.init();
});
