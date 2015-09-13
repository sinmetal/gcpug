/// <reference path="../../typings/jquery/jquery.d.ts" />
/// <reference path="../../typings/angularjs/angular.d.ts" />
/// <reference path="../../typings/materialize/materialize.d.ts" />

module Gcpug {

    export class EventController {

        constructor(private $scope:ng.IScope, private $http:ng.IHttpService) {
            var _this = this;
            _this.items = [];
            _this.getEvents(10);
        }

        getEvents(limit: number) {
            if (!limit) {
                limit = 10;
            }
            var _this = this;
            _this.$http.get('/api/1/event?limit=' + limit)
                .success(function (data, status, headers, config) {
                    for (var key in data) {
                        data[key].status = Gcpug.EventFilter.getColorClassByDate(data[key].startAt);
                    }
                    _this.items = data;
                });
        }

        items:any[];

    }
}