/// <reference path="../typings/jquery/jquery.d.ts" />
/// <reference path="../typings/angularjs/angular.d.ts" />
/// <reference path="../typings/materialize/materialize.d.ts" />
/// <reference path="../typings/moment-timezone/moment-timezone.d.ts" />

module Gcpug {

	export class EventController {
		constructor(private $scope : ng.IScope, private $http : ng.IHttpService) {
			var _this = this;
				$http.get('/api/1/event?limit=5')
				.success(function(data, status, headers, config) {
					for(var key in data) {
						data[key].status = Gcpug.Filter.getColorClassByDate(data[key].startAt);
					}
					_this.items = data;
				});

		}
		items : {};
	}

	export class OrganizationController {
		constructor(private $scope : ng.IScope, private $http : ng.IHttpService) {
			var _this = this;
			$http.get('/api/1/organization')
				.success(function(data, status, headers, config) {
					_this.items = data;
				});
		}
		items : {};
	}

	export interface TopInterface extends ng.IModule {}

	export class Filter {

		static getColorClassByDate(time : string) {
			var startAt = moment(time);
			if ( startAt.isSame(moment(), 'day') ) {
				return 'amber lighten-5';
			}
			if ( startAt.isBefore(moment(), 'day') ) {
				return 'grey lighten-3';
			}
			return '';
		}

		static formatDatetime(time : string) {
			var startAt = moment(time);
			return startAt.format('YYYY/M/D H:mm')+'〜';
		}
	}
}

var app: Gcpug.TopInterface = angular.module('Gcpug', []);

app.controller('EventController', ['$scope', '$http', function($scope, $http) {
	return new Gcpug.EventController($scope, $http);
}]);

app.controller('OrganizationController', ['$scope', '$http', function($scope, $http) {
	return new Gcpug.OrganizationController($scope, $http);
}]);

app.filter('formatDatetime', function() {
	return Gcpug.Filter.formatDatetime;
});