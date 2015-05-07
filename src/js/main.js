var Gcpug;
(function (Gcpug) {
    var EventController = (function () {
        function EventController($scope, $http) {
            this.$scope = $scope;
            this.$http = $http;
            var _this = this;
            $http.get('/api/1/event?limit=5').success(function (data, status, headers, config) {
                for (var key in data) {
                    data[key].status = Gcpug.EventFilter.getColorClassByDate(data[key].startAt);
                }
                _this.items = data;
            });
        }
        return EventController;
    })();
    Gcpug.EventController = EventController;
})(Gcpug || (Gcpug = {}));
var Gcpug;
(function (Gcpug) {
    var OrganizationController = (function () {
        function OrganizationController($scope, $http) {
            this.$scope = $scope;
            this.$http = $http;
            var _this = this;
            $http.get('/api/1/organization').success(function (data, status, headers, config) {
                _this.items = data;
            });
        }
        return OrganizationController;
    })();
    Gcpug.OrganizationController = OrganizationController;
})(Gcpug || (Gcpug = {}));
var Gcpug;
(function (Gcpug) {
    var EventFilter = (function () {
        function EventFilter() {
        }
        EventFilter.getColorClassByDate = function (time) {
            var startAt = moment(time);
            if (startAt.isSame(moment(), 'day')) {
                return 'amber lighten-5';
            }
            if (startAt.isBefore(moment(), 'day')) {
                return 'grey lighten-3';
            }
            return '';
        };
        EventFilter.formatDatetime = function (time) {
            var startAt = moment(time);
            return startAt.format('YYYY/M/D H:mm') + '〜';
        };
        return EventFilter;
    })();
    Gcpug.EventFilter = EventFilter;
})(Gcpug || (Gcpug = {}));
var app = angular.module('Gcpug', []);
app.controller('EventController', ['$scope', '$http', function ($scope, $http) {
    return new Gcpug.EventController($scope, $http);
}]);
app.controller('OrganizationController', ['$scope', '$http', function ($scope, $http) {
    return new Gcpug.OrganizationController($scope, $http);
}]);
app.filter('formatDatetime', function () {
    return Gcpug.EventFilter.formatDatetime;
});
