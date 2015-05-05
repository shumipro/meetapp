import $ from 'jquery'
import constants from '../constants'
import util from '../util'
import Handlebars from 'handlebars'

var html =  '{{#if hasPrev}}<span class="prev">' +
                '<a href="/app/list?page={{prev}}">前へ</a>' +
            '</span>{{/if}}' + 
            '{{#each pages}}{{#if isCurrent}}<span class="page">{{page}}</span>{{else}}<span class="page"><a href="/app/list?page={{page}}">{{page}}</a></span>{{/if}}{{/each}}' +
            '{{#if hasNext}}<span class="next">' +
               '<a href="/app/list?page={{next}}">次へ</a>' +
            '</span>{{/if}}'

var tmpl = Handlebars.compile(html)

export default class ConstantFilter {
    constructor() {
        var $wrap = $('.ma-app-list-pagination'),
            $container = $wrap.find('.pagination'),
            currentPage = $wrap.data('current-page'),
            perPage = $wrap.data('per-page'),
            totalCount = $wrap.data('total-count')

        var pages = []
        // temp
        var min = 1,
            max = Math.floor((totalCount - 1) / perPage) + 1
        for(var i=min; i<=max; i++){
            var obj = {page: i}
            obj.isCurrent = (i === currentPage - 0)
            pages.push(obj)
        }
        $(tmpl({
            currentPage: currentPage,
            perPage: perPage,
            totalCount: totalCount,
            pages: pages,
            hasPrev: currentPage > 1,
            hasNext: currentPage < max,
            prev: currentPage - 1,
            next: currentPage + 1
        })).appendTo($container)
    }
}