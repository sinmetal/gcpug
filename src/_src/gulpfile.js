var gulp = require('gulp');
var compass = require('gulp-compass');
var typescript = require('gulp-typescript');
var uglify = require('gulp-uglify');
var plumber = require('gulp-plumber');

gulp.task('typescript', function() {
	return gulp.src('typescript/**/**.ts')
		.pipe(plumber())
		.pipe(typescript({
			out : 'main.js',
			removeComments : true
		}))
		//.pipe(uglify())
		.pipe(gulp.dest('../js'));
});

gulp.task('compass', function() {
	gulp.src('sass/*.scss')
		.pipe(plumber())
		.pipe(compass({
			config_file: './config.rb',
			css: '../css',
			sass: 'sass',
			comments: false
		}));
	gulp.src('sass/materialize/*.scss')
		.pipe(plumber())
		.pipe(compass({
			config_file: './config.rb',
			css: '../css',
			sass: 'sass/materialize',
			comments: false
		}));
});

gulp.task('watch', [ 'typescript',  'compass'], function() {
	 gulp.watch('typescript/**', ['typescript']);
	gulp.watch('sass/**', ['compass']);
});