/// <reference path="../../typings/jquery/jquery.d.ts" />
/// <reference path="../../typings/angularjs/angular.d.ts" />
/// <reference path="../../typings/materialize/materialize.d.ts" />

module Gcpug {

    export class EventFilter {

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
            var zone = moment().utcOffset();
            return startAt.utcOffset(zone).format('YYYY/M/D H:mm')+'〜';
        }
    }
}