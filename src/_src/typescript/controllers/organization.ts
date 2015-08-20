/// <reference path="../../typings/jquery/jquery.d.ts" />
/// <reference path="../../typings/angularjs/angular.d.ts" />
/// <reference path="../../typings/materialize/materialize.d.ts" />

module Gcpug {

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
}