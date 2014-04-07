angular.module('ui.keyvalue', [])
    .controller('KeyValueCtrl', ['$scope', '$attrs',
        function($scope, $attrs) {
            $scope.$watch('keyvalue', function(value) {
                if (value) {
                    $scope.values = value;
                }
            });

            // $scope.values = [];
            $scope.values.value = function() {
                var tmp = {};
                $.each(this, function(index, param) {
                    if (param.active) {
                        tmp[param.key] = param.value;
                    }
                });
                return tmp;
            };

            $scope.addValue = function() {
                // Set the last row to be active
                if ($scope.values.length > 0) {
                    $scope.values[$scope.values.length - 1].active = true;
                }
                $scope.values.push({
                    active: false
                });
            };
            $scope.deleteValue = function(index) {
                if (!$scope.values[index].preset) {
                    $scope.values.splice(index, 1);
                }
            };
        }
    ])
    .directive('keyvalue', function() {
        return {
            restrict: 'EA',
            transclude: false,
            replace: true,
            scope: {
                values: '=ngModel'
            },
            templateUrl: '/js/templates/keyvalue/keyvalue.html',
            controller: 'KeyValueCtrl',
            link: function(scope, elem, attrs) {
                // Ensure that the values array has an inactive element
                if (scope.values.length === 0 || scope.values[scope.values.length - 1].active) {
                    scope.values.push({
                        active: false
                    });
                }
            }
        };
    });
