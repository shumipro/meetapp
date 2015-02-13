var gulp          = require('gulp');
var shell         = require('gulp-shell');
var runSequence   = require('run-sequence');
var jshint        = require('gulp-jshint');
var webpack       = require('gulp-webpack');
var stylus        = require('gulp-stylus');
var pleeease      = require('gulp-pleeease');
var minifyCss     = require('gulp-minify-css');
var rename        = require('gulp-rename');
var webserver     = require('gulp-webserver');


// :stylus
gulp.task('stylus', function(){
  return gulp.src([
    './public/stylus/*.styl'
  ])
  .pipe(stylus())
  .pipe(pleeease())
  .pipe(minifyCss({keepSpecialComments: 0}))
  .pipe(rename({extname: '.css'}))
  .pipe(gulp.dest('./public/css/'));
});

// :Livereload for Static Build
gulp.task('webserver', function() {
  gulp.src('public/')
  .pipe(webserver({
    livereload: true,
    port: 8088,
    // directoryListing: true,
    open: true
  }));
});


gulp.task('jshint', function() {
    return gulp.src('public/js/**/*.js')
        .pipe(jshint({
          esnext: true,
          asi: true
        }))
        .pipe(jshint.reporter('default'));
});

gulp.task('watch', function() {
    gulp.watch(['public/js/**/*.js'], ['jshint', 'webpack']);
    gulp.watch(['public/css/stylus/**/*.styl'],['stylus']);
});

gulp.task('webpack', function() {
   gulp.src('public/js/ma/main.js')
    .pipe(webpack({
      output: {
        filename: 'public/dist/main.js'
      },
      devtool: 'inline-source-map',
      resolve: {
        extensions: ['', '.js', '.jsx']
      },
      module: {
        loaders: [
          // { test: /\.jsx$/, loader: 'jsx-loader?harmony' }
          { test: /\.js|\.jsx$/, exclude: /node_modules|public\/dist/, loader: '6to5-loader' }
        ]
      }
    }))
    .pipe(gulp.dest(''));
});

gulp.task('default', function(){
    runSequence('webpack', 'watch', 'stylus', 'webserver');
});

gulp.task('server', shell.task(['npm start']));

gulp.task('test', shell.task(['npm test']));