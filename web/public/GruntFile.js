module.exports = function(grunt) {

    // Project configuration.
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        clean: {
            css: ['css/combined.css', 'css/source-map.css'],
            js: ['js/bootstrap.js', 'js/combined.js', 'js/combined.min.js', 'js/combined.map', 'combined.min.map']
        },
        recess: {
            lint: {
                options: {
                    noOverqualifying: false
                },
                src: [
                    'less/style.less'
                ]
            },
        },
        less: {
            development: {
                options: {
                    compress: false,
                    sourceMap: true,
                    // sourceMapFilename: 'css/source-map.css',
                    // sourceMapURL: 'source-map.css',
                    sourceMapRootpath: '/'
                },
                files: {
                    'css/combined.css': ['less/bootstrap.less', 'less/style.less']
                }
            },
            build: {
                options: {
                    compress: true
                },
                files: {
                    'css/combined.min.css': ['less/bootstrap.less', 'less/style.less']
                }
            }
        },
        jshint: {
            gruntfile: ['GruntFile.js'],
            development: ['js/**/*.js', '!js/vendors/**/*.js', '!js/bower_components/**/*.js', '!js/bootstrap.js', '!js/bootstrap/**/*.js', '!js/source-map.js'],
            build: ['js/**/*.js', '!js/vendors/**/*.js', '!js/bower_components/**/*.js', '!js/bower_components/**/*.js', '!js/bootstrap/**/*.js', '!js/source-map.js']
        },
        jsbeautifier: {
            files: ['js/**/*.js', '!js/combined.js', '!js/combined.min.js', '!js/bootstrap.js', '!js/source-map.js']
        },
        concat: {
            options: {
                separator: ';',
            },
            js: {
                src: [
                    'js/vendors/jquery.js',
                    'js/vendors/holder.js',
                    'js/vendors/jquery.actual.js',
                    'js/vendors/jquery.selector.js',
                    'js/vendors/markdown.js',
                    'js/bootstrap/transition.js',
                    'js/bootstrap/alert.js',
                    'js/bootstrap/button.js',
                    'js/bootstrap/carousel.js',
                    'js/bootstrap/collapse.js',
                    'js/bootstrap/dropdown.js',
                    'js/bootstrap/modal.js',
                    'js/bootstrap/tooltip.js',
                    'js/bootstrap/popover.js',
                    'js/bootstrap/crollspy.js',
                    'js/bootstrap/tab.js',
                    'js/bootstrap/affix.js',
                    'js/bootstrap/markdown.js',
                    'bower_components/angular/angular.js',
                    'js/vendors/underscore.js',
                    'js/vendors/backbone.js',
                    'js/vendors/pretty-json.js',
                    'bower_components/ace-builds/src-min-noconflict/ace.js',
                    'bower_components/angular-ui-ace/ui-ace.js',
                    'bower_components/ace-builds/src-min-noconflict/mode-javascript.js',
                    'bower_components/angular-bootstrap/ui-bootstrap-tpls.js',
                    'js/src/**/*.js',
                ],
                dest: 'js/combined.js',
            },
        },
        uglify: {
            development: {
                options: {
                    beautify: true,
                    compress: false,
                    preserveComments: 'all',
                    dead_code: true,
                    mangle: false
                },
                src: 'js/combined.js',
                dest: 'js/combined.js',
            },
            build: {
                options: {
                    banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n',
                    sourceMap: 'js/source-map.js',
                    dead_code: true,
                    mangle: false
                },
                src: 'js/combined.js',
                dest: 'js/combined.min.js'
            }
        },
        watch: {
            gruntfile: {
                files: 'Gruntfile.js',
                tasks: ['jshint:gruntfile'],
            },
            development_html: {
                files: ['**/*.html'],
                options: {
                    // Start a live reload server on the default port 35729
                    livereload: true,
                },
            },
            development_css: {
                files: ['less/**/*.less'],
                tasks: ['clean:css', 'recess:lint', 'less:development'],
                options: {
                    // Start a live reload server on the default port 35729
                    livereload: true,
                },
            },
            development_js: {
                files: ['js/**/*.js', '!js/combined.js', '!js/combined.min.js', '!js/bootstrap.js', '!js/source-map.js'],
                tasks: ['clean:js', 'jshint:development', 'jsbeautifier', 'concat:js', 'uglify:development'],
                options: {
                    // Start a live reload server on the default port 35729
                    livereload: true,
                },
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-contrib-jshint');
    grunt.loadNpmTasks('grunt-contrib-less');
    grunt.loadNpmTasks('grunt-jsbeautifier');
    grunt.loadNpmTasks('grunt-recess');

    // Default tasks.
    grunt.registerTask('default', ['build_dev']);
    grunt.registerTask('dev', ['build_dev', 'watch']);
    grunt.registerTask('build_dev', ['clean:css', 'clean:js', 'recess:lint', 'less:development', 'jshint:development', 'jsbeautifier', 'concat:js', 'uglify:development']);
    grunt.registerTask('build_prod', ['clean:css', 'clean:js', 'recess:lint', 'less:build', 'jshint:build', 'concat:js', 'uglify:build']);

};
