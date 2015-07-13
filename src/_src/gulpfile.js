var gulp = require('gulp');
var sass = require('gulp-sass');
var typescript = require('gulp-typescript');
var plumber = require('gulp-plumber');
var bourbon = require('node-bourbon');

gulp.task('typescript', function() {
	return gulp.src('typescript/**/**.ts')
		.pipe(plumber())
		.pipe(typescript({
			out : 'main.js',
			removeComments : true
		}))
		.pipe(gulp.dest('../js'));
});

gulp.task('sass', function() {
	gulp.src('sass/materialize/*.scss')
		.pipe(sass({
			includePaths: bourbon.includePaths
		}))
		.pipe(plumber())
		.pipe(gulp.dest('../css'));
})

gulp.task('watch', [ 'typescript',  'sass'], function() {
	 gulp.watch('typescript/**', ['typescript']);
	gulp.watch('sass/**', ['sass']);
});