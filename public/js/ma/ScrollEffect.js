import $ from 'jquery'
import waypoint from 'waypoints/lib/jquery.waypoints'

export default class ScrollEffect {
    constructor() {
        this.smoothScroll()
        this.animateScript()
        this.scrollMenu()
    }

    smoothScroll() {
        $('a[href*=#]:not([href=#carousel-example-generic])').click(function() {
            if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
                var target = $(this.hash);
                target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
                if (target.length) {
                    $('html,body').animate({
                        scrollTop: target.offset().top
                    }, 1000);
                    return false;
                }
            }
        });
    }

    animateScript() {
        $('.scrollpoint.sp-effect1').waypoint(function() {
            $(this.element).toggleClass('active');
            $(this.element).toggleClass('animated fadeInLeft');
        }, {
            offset: '100%'
        });
        $('.scrollpoint.sp-effect2').waypoint(function() {
            $(this.element).toggleClass('active');
            $(this.element).toggleClass('animated fadeInRight');
        }, {
            offset: '100%'
        });
        $('.scrollpoint.sp-effect3').waypoint(function() {
            $(this.element).toggleClass('active');
            $(this.element).toggleClass('animated fadeInDown');
        }, {
            offset: '100%'
        });
        $('.scrollpoint.sp-effect4').waypoint(function() {
            $(this.element).toggleClass('active');
            $(this.element).toggleClass('animated fadeIn');
        }, {
            offset: '100%'
        });
        $('.scrollpoint.sp-effect5').waypoint(function() {
            $(this.element).toggleClass('active');
            $(this.element).toggleClass('animated fadeInUp');
        }, {
            offset: '100%'
        });
    }

    scrollMenu() {
        var num = 50;
        if ($(window).scrollTop() > num) {
            $('nav').addClass('scrolled');
        }
        $(window).bind('scroll', function() {
            if ($(window).scrollTop() > num) {
                $('nav').addClass('scrolled');
            } else {
                $('nav').removeClass('scrolled');
            }
        });
        $('ul.navbar-nav li a').bind('click', function() {
            if ($(this).closest('.navbar-collapse').hasClass('in')) {
                $(this).closest('.navbar-collapse').removeClass('in');
            }
        });
    }
}