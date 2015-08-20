/// <reference path="../typings/jquery/jquery.d.ts" />
/// <reference path="../typings/angularjs/angular.d.ts" />
/// <reference path="../typings/materialize/materialize.d.ts" />
/// <reference path="../typings/moment-timezone/moment-timezone.d.ts" />

/// <reference path="controllers/event.ts" />
/// <reference path="controllers/organization.ts" />
/// <reference path="filters/eventFilter.ts" />

var app: ng.IModule = angular.module('Gcpug', []);

app.controller('EventController', ['$scope', '$http', function($scope, $http) {
	return new Gcpug.EventController($scope, $http);
}]);

app.controller('OrganizationController', ['$scope', '$http', function($scope, $http) {
	return new Gcpug.OrganizationController($scope, $http);
}]);

app.filter('formatDatetime', function() {
	return Gcpug.EventFilter.formatDatetime;
});