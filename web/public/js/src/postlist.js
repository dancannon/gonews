$(function() {
    // The height of the content block when it's not expanded
    var adjustheight = 125;
    // The "more" link text
    var moreText = "+  More";
    // The "less" link text
    var lessText = "- Less";

    if ($(".text-post-type .post-details .post-description").actual('height') > adjustheight) {
        // Sets the .post-description div to the specified height and hides any content that overflows
        $(".text-post-type .post-details .post-description").css('height', adjustheight).css('overflow', 'hidden');

        // The section added to the bottom of the "post-details" div
        $(".text-post-type .post-details").append('<a href="#" class="adjust"></a>');
        // Set the "More" text
        $("a.adjust").text(moreText);

        $(".adjust").click(function(e) {
            if ($(this).parents("div.post-details").find(".post-description").data('expanded') !== true) {
                $(this).parents("div.post-details").find(".post-description").css('height', 'auto').css('overflow', 'visible').data('expanded', true);
                $(this).text(lessText);
            } else {
                $(this).parents("div.post-details").find(".post-description").css('height', adjustheight).css('overflow', 'hidden').data('expanded', false);
                $(this).text(moreText);
            }

            e.preventDefault();
        });
    }
});
