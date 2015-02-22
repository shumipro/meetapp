var express = require('express');
var router = express.Router();

var SORT_LABELS = {
    'new': {
        title: '新着アプリ'
    },
    'popular': {
        title: '人気アプリ'
    }
};

router.get('/list', function(req, res, next) {
    var orderBy = req.query["orderBy"];
    // TODO: retrieve from DB

    var label = SORT_LABELS[orderBy].title;
    res.render('app/list', {
        title: 'MeetApp - ' + label,
        subTitle: label
    });
});

router.get('/register', function(req, res, next) {
    res.render('app/register', {
        title: 'MeetApp - アプリの登録'
    });
});

module.exports = router;