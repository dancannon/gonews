// define "models"
function Value(props) {
    props = props || {};
    props.name = props.name || '';
    props.type = props.type || 'value';
    props.value = props.value || '';
    props.params = props.params || [];

    this.name = props.name;
    this.type = props.type;
    this.value = props.value;
    this.params = props.params;
}


var ruleBuilder = angular.module('ruleBuilder', ['ui.keyvalue', 'ui.ace', 'ui.bootstrap']);

ruleBuilder.config(function($locationProvider) {
    $locationProvider.html5Mode(true); // enable pushState
});

ruleBuilder.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});

ruleBuilder.filter('empty', function() {
    return function(input) {
        return _.isEmpty(input);
    };
});

ruleBuilder.controller('RuleCtrl', function RuleCtrl($scope, $http, $location, $sce) {
    $scope.extractors = [{
        id: "oembed",
        name: "OEmbed",
        params: [{
            key: "endpoint",
            name: "OEmbed Endpoint (Required)",
            value: "",
            required: true,
            preset: true
        }]
    }, {
        id: "opengraph",
        name: "OpenGraph",
        params: []
    }, {
        id: "selector",
        name: "CSS Selector (Raw Content)",
        params: [{
            key: "selector",
            name: "CSS Selector (Required)",
            value: "",
            required: true,
            preset: true
        }, {
            key: "attribute",
            name: "Attribute",
            value: "",
            preset: true
        }, {
            key: "restype",
            name: "Result Type",
            value: "",
            choices: {
                "": "",
                "first": "First Result",
                "all": "All Results",
                "merge": "Merge Results"
            },
            preset: true
        }]
    }, {
        id: "selector_text",
        name: "CSS Selector (Text Only)",
        params: [{
            key: "selector",
            name: "CSS Selector (Required)",
            value: "",
            required: true,
            preset: true
        }, {
            key: "attribute",
            name: "Attribute",
            value: "",
            preset: true
        }, {
            key: "restype",
            name: "Result Type",
            value: "",
            choices: {
                "": "",
                "first": "First Result",
                "all": "All Results",
                "merge": "Merge Results"
            },
            preset: true
        }]
    }, {
        id: "text",
        name: "Text",
        params: [{
            key: "format",
            name: "Result Format",
            value: "raw",
            choices: {
                "raw": "Raw",
                "text": "Text",
                // "markdown": "Markdown"
            },
            preset: true
        }]
    }, {
        id: "title",
        name: "Page Title",
        params: []
    }, {
        id: "url_mapper",
        name: "URL mapper",
        params: [{
            key: "pattern",
            name: "Pattern",
            value: "",
            required: true,
            preset: true
        }, {
            key: "replacement",
            name: "Replacement",
            value: "",
            required: true,
            preset: true
        }]
    }, {
        id: "javascript",
        name: "JS Script",
        params: [{
            key: "script",
            name: "Script",
            value: "",
            required: true,
            preset: true
        }]
    }];
    $scope.types = {
        "extractor": [
            new Value({
                type: 'extractor'
            })
        ],
        "general": [
            new Value({
                name: 'title',
                type: 'extractor',
                value: 'title',
            }),
            new Value({
                name: 'content'
            })
        ],
        "text": [
            new Value({
                name: 'title',
                type: 'extractor',
                value: 'title'
            }),
            new Value({
                name: 'text',
                type: 'extractor',
                value: 'text'
            })
        ],
        "image": [
            new Value({
                name: 'title',
                type: 'extractor',
                value: 'title',
            }),
            new Value({
                name: 'caption',
                type: 'extractor',
                value: 'selector',
            }),
            new Value({
                name: 'author',
                type: 'values',
                value: [
                    new Value({
                        name: 'name',
                        type: 'extractor',
                        value: 'selector',
                    }),
                    new Value({
                        name: 'url',
                        type: 'extractor',
                        value: 'selector',
                    }),
                ],
            }),
            new Value({
                name: 'thumbnail',
                type: 'values',
                value: [
                    new Value({
                        name: 'url',
                        type: 'extractor',
                        value: 'selector',
                    }),
                    new Value({
                        name: 'width',
                        type: 'extractor',
                        value: 'selector',
                    }),
                    new Value({
                        name: 'height',
                        type: 'extractor',
                        value: 'selector',
                    }),
                ],
            }),
            new Value({
                name: 'url',
                type: 'extractor',
                value: 'selector',
            }),
            new Value({
                name: 'height',
                type: 'extractor',
                value: 'selector',
            }),
            new Value({
                name: 'width',
                type: 'extractor',
                value: 'selector',
            })
        ],
        "video": [
            new Value({
                name: 'title',
                type: 'extractor',
                value: 'title',
            }),
            new Value({
                name: 'description',
                type: 'extractor',
                value: 'selector',
            }),
            new Value({
                name: 'author',
                type: 'values',
                value: [
                    new Value({
                        name: 'name',
                        type: 'extractor',
                        value: 'selector',
                    }),
                    new Value({
                        name: 'url',
                        type: 'extractor',
                        value: 'selector',
                    }),
                ],
            }),
            new Value({
                name: 'thumbnail',
                type: 'values',
                value: [
                    new Value({
                        name: 'url',
                        type: 'extractor',
                        value: 'selector',
                    }),
                    new Value({
                        name: 'width',
                        type: 'extractor',
                        value: 'selector',
                    }),
                    new Value({
                        name: 'height',
                        type: 'extractor',
                        value: 'selector',
                    }),
                ],
            }),
            new Value({
                name: 'html',
                type: 'extractor',
                value: 'selector',
            }),
        ],
        "rich": [
            new Value({
                name: 'title',
                type: 'extractor',
                value: 'title',
            }),
            new Value({
                name: 'html',
                type: 'extractor',
                value: 'selector',
            }),
            new Value({
                name: 'height',
                type: 'extractor',
                value: 'selector',
            }),
            new Value({
                name: 'width',
                type: 'extractor',
                value: 'selector',
            })
        ]
    };

    // Form utility fields
    $scope.error = '';
    $scope.action = '';
    $scope.submitted = false;
    $scope.mode = "simple";

    // Rule fields
    $scope.id = "";
    $scope.name = "";
    $scope.test_url = "";
    $scope.host = "";
    $scope.path_pattern = "/.*";
    $scope.type = "unknown";
    $scope.values = [];

    // New value scope fields
    $scope.container = {
        newname: '',
        selected: null,
        restype: "",
    };

    // Load rule if rule ID is set
    if (!_.isEmpty(rule)) {
        $scope.id = rule.id || "";
        $scope.name = rule.name || "";
        $scope.host = rule.host || "";
        $scope.path_pattern = rule.path_pattern || "/.*";
        $scope.test_url = rule.test_url || "";
        $scope.type = rule.type || "unknown";

        var loadValues = function(values) {
            var tmp = [];

            $.each(values, function(index, value) {
                if (value.type === "extractor") {
                    var params = getExtractorParams(value.id, false);
                    $.each(value.params, function(pk, pv) {
                        var found = false;
                        $.each(params, function(pindex, defaultParam) {
                            if (defaultParam.key === pk) {
                                params[pindex] = _.extend(_.clone(defaultParam), {
                                    value: pv
                                });
                                found = true;
                            }
                        });

                        if (!found) {
                            params.push({
                                key: pk,
                                value: pv,
                                active: true,
                            });
                        }
                    });

                    params.push({
                        active: false
                    });

                    tmp.push(new Value({
                        name: value.name,
                        type: "extractor",
                        value: value.id,
                        params: params,
                    }));
                } else {
                    if (value.type === 'values') {
                        var childValues = loadValues(value);
                        tmp.push(new Value({
                            name: value.name,
                            type: "values",
                            value: loadValues(value.value),
                        }));
                    } else if (value.type === 'value') {
                        tmp.push(new Value(value));
                    }
                }
            });

            return tmp;
        };

        if (rule.values !== null) {
            $scope.values = loadValues(rule.values);
        }
    }

    // Build values object
    var prepareValues = function(values) {
        var tmp;
        if (typeof values === 'object') {
            if (values instanceof Value) {
                if (values.type === 'extractor') {
                    return {
                        "name": values.name,
                        "id": values.value,
                        "type": "extractor",
                        "params": _.object(_.map(values.params, function(x) {
                            return [x.key, x.value];
                        }))
                    };
                } else {
                    return {
                        "name": values.name,
                        "type": values.type,
                        "value": prepareValues(values.value),
                    };
                }
            } else {
                tmp = [];

                $.each(values, function(index, value) {
                    tmp.push(prepareValues(value));
                });

                return tmp;
            }
        } else if (Array.isArray(values)) {
            tmp = [];

            $.each(values, function(index, value) {
                tmp.push(prepareValues(value));
            });

            return tmp;
        } else {
            return values;
        }
    };

    $scope.onSubmit = function(event) {
        $scope.submitted = true;

        if ($scope.form.$invalid) {
            return;
        }

        if ($scope.action == 'test') {
            $scope.onSubmitTestRule();
        } else {
            $scope.onSubmitSaveRule();
        }
    };
    $scope.onSubmitSaveRule = function() {
        $http.post('/rule/save', {
            "post_id": post_id,
            "rule": {
                "id": $scope.id,
                "name": $scope.name,
                "type": $scope.type === 'extractor' ? 'unknown' : $scope.type,
                "host": $scope.host,
                "path_pattern": $scope.path_pattern,
                "test_url": $scope.test_url,
                "values": prepareValues($scope.values),
            }
        }).success(function(data, status, headers, config) {
            $scope.error = "";
            $scope.id = data.id;
            if ($location.path() !== "/rule/edit/" + data.id) {
                $location.path("/rule/edit/" + data.id);
            }

            var $el = $("#result-preview");
            var node = new PrettyJSON.view.Node({
                el: $el,
                data: data.response
            });
            window.location.hash = '#preview';
        }).error(function(data, status, headers, config) {
            if (data.error) {
                $scope.error = "That rule was not valid (" + data.error + ")";
            } else {
                if (status == 400) {
                    $scope.error = "That rule was not valid";
                } else {
                    $scope.error = "An error occurred, please try again later";
                }
            }
        });
    };
    $scope.onSubmitTestRule = function() {
        $http.post('/rule/test', {
            "rule": {
                "name": $scope.name,
                "type": $scope.type === 'extractor' ? 'unknown' : $scope.type,
                "host": $scope.host,
                "path_pattern": $scope.path_pattern,
                "test_url": $scope.test_url,
                "values": prepareValues($scope.values),
            }
        }).success(function(data, status, headers, config) {
            var $el = $("#result-preview");
            var node = new PrettyJSON.view.Node({
                el: $el,
                data: data.response
            });
            window.location.hash = '#preview';
        }).error(function(data, status, headers, config) {
            if (data.error) {
                $scope.error = "That rule was not valid (" + data.error + ")";
            } else {
                if (status == 400) {
                    $scope.error = "That rule was not valid";
                } else {
                    $scope.error = "An error occurred, please try again later";
                }
            }
        });
    };

    function getExtractorParams(name, addRow) {
        if (typeof addRow === "undefined") {
            addRow = true;
        }

        var params = [];
        $.each($scope.extractors, function(index, extractor) {
            if (name === extractor.id) {
                params = _.extend([], extractor.params);
                params = _.map(params, function(param) {
                    param.active = true;
                    return param;
                });

                if (addRow) {
                    params.push({
                        active: false
                    });
                }

                return false;
            }
        });

        return params;
    }

    $scope.onUpdateExtractorType = function(value) {
        value.params = getExtractorParams(value.value);
    };

    $scope.onUpdateRuleType = function() {
        $.each($scope.types, function(id, type) {
            if ($scope.type == id) {
                $scope.values = type;
                $.each($scope.values, function(key, value) {
                    if (value.type === 'extractor') {
                        value.params = getExtractorParams(value.value);
                    }
                });

                return;
            }
        });
    };

    $scope.addValue = function() {
        if ($scope.container.newname !== '') {
            $scope.values = $scope._addValue($scope.values, new Value({
                name: $scope.container.newname,
            }));
            $scope.container.newname = "";
        }
    };
    $scope.addChildValue = function(parent) {
        if (parent.newname !== '') {
            parent.value = $scope._addValue(new Value({
                name: parent.newname
            }));
            parent.newname = '';
        }
    };
    $scope.addSimpleValue = function() {
        if ($scope.container.newname !== '') {
            // Get selectors
            var elements = $("#rule-builder-simple iframe").contents().find(".highlight-selected");
            elements.removeClass("highlight-selected");
            var selectors = _.map(elements, function(el) {
                return $(el).getSelector({
                    ignore: {
                        classes: ['highlight-selected', 'highlight']
                    }
                });
            });

            $scope.values = $scope._addValue($scope.values, new Value({
                name: $scope.container.newname,
                value: "selector",
                type: "extractor",
                params: [{
                    active: true,
                    key: "selector",
                    name: "CSS Selector (Required)",
                    preset: true,
                    required: true,
                    value: selectors.join(", ")
                }, {
                    active: true,
                    key: "restype",
                    name: "Result Type",
                    preset: true,
                    required: false,
                    value: $scope.container.restype,
                }, {
                    active: false
                }]
            }));
            $scope.container.newname = "";
        }
    };
    $scope._addValue = function(values, value) {
        if (!Array.isArray(values)) {
            values = [];
        }

        for (var i = values.length - 1; i >= 0; i--) {
            if (values[i].name == value.name) {
                values[i] = value;
                return values;
            }
        }

        values.push(value);

        return values;
    };
    $scope.updateValue = function() {
        if ($scope.container.selected !== null) {
            var iframe = $("#rule-builder-simple-iframe").contents();

            // Get selectors
            var elements = iframe.find(".highlight-selected");
            elements.removeClass("highlight-selected");
            var selectors = _.map(elements, function(el) {
                return $(el).getSelector({
                    ignore: {
                        classes: ['highlight-selected', 'highlight']
                    }
                });
            });
            var selector = selectors.join(", ");

            // Update selector value
            for (var i = $scope.container.selected.params.length - 1; i >= 0; i--) {
                var param = $scope.container.selected.params[i];
                if (param.key === "selector") {
                    param.value = selector;
                }
            }

            // Readd the highlights
            iframe.find(selector).addClass("highlight-selected");
        }
    };
    $scope.deleteValue = function(index) {
        $scope.values.splice(index, 1);
    };
    $scope.onUpdateValueType = function(value) {
        if (value.type == 'value') {
            value.value = '';
        } else if (value.type == 'values') {
            value.value = [];
        } else if (value.type == 'extractor') {
            value.value = '';
        }
    };
    $scope.onUpdateSelectedValues = function() {
        var iframe = $("#rule-builder-simple-iframe").contents();
        iframe.find(".highlight-selected").removeClass("highlight-selected");
        if ($scope.container.selected !== null) {
            // get selector
            var selector = "";
            for (var i = $scope.container.selected.params.length - 1; i >= 0; i--) {
                var param = $scope.container.selected.params[i];
                if (param.key === "selector") {
                    selector = param.value;
                }
            }

            iframe.find(selector).addClass("highlight-selected");
        }
    };
    $scope.onTestUrlChange = function() {
        if ($scope.mode === "simple") {
            $scope.initSimple();
        }
    };

    $scope.trustSrc = function(src) {
        return $sce.trustAsResourceUrl(src);
    };

    $scope.initSimple = function() {
        window.updateTestUrl = function(test_url) {
            $scope.test_url = test_url;
        };

        $scope.mode = "simple";
    };
    $scope.initAdvanced = function() {
        $scope.mode = "advanced";
    };
    $scope.startIframe = function() {
        $("#rule-builder-simple-iframe")
            .attr("src", "/rule/load_url?url=" + $scope.test_url)
            .removeClass("hide");
        $("#rule-builder-simple .navbar .navbar-form").removeClass("hide");
        $("#rule-builder-simple-placeholder").addClass("hide");
    };
    $scope.reloadIframe = function() {
        $("#rule-builder-simple-iframe").attr("src", "/rule/load_url?url=" + $scope.test_url);
    };

    // Init rule builder
    $scope.initSimple();
});

var $el = $("#result-preview");
var node = new PrettyJSON.view.Node({
    el: $el,
    data: {},
});
