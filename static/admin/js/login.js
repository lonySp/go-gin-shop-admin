var app = {
    init: function() {
        this.getCaptcha();
        this.captchaImgChange();
    },
    getCaptcha: function() {
        $.get("/admin/captcha?t=" + Math.random(), function(response) {
            console.log(response);
            console.log("特殊测撒打算阿斗撒的阿萨德")
            $("#captchaId").val(response.captchaId);
            $("#captchaImg").attr("src", response.captchaImg);
            $("#username").val("admin");
            $("#password").val("123456");
            $("#verify").val(response.answer);
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
