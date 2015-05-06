/// <reference path="../typings/jquery/jquery.d.ts" />
/// <reference path="../typings/materialize/materialize.d.ts" />
/// <reference path="../typings/handlebars/handlebars.d.ts" />
/// <reference path="../typings/moment-timezone/moment-timezone.d.ts" />
/// <reference path="../typings/es6-promise/es6-promise.d.ts" />
module GCPUG {

	export class Api {
		getTemplate(name : string) {
			return new Promise(function(resolve) {
				if ($('#template-' + name).length === 0) {
					$.ajax({
						method : 'get',
						url : '/template/' + name + '.html',
						dataType : 'text',
						success : function (data, status, jqXHR) {
							$('body').append(data);
							resolve('success');
						}
					});
				} else {
					resolve('success');
				}
			});
		}

		getEventList() {
			var getTemplate = this.getTemplate('event-list');
			getTemplate.then(function() {
				$.ajax({
					method : 'get',
					url : '/api/1/event?limit=5',
					success : function (data, status, jqXHR) {
						var source = $('#template-event-list').html();
						var template = Handlebars.compile(source);
						var html = template({list:data});
						$('#event-list').append(html);
					}
				});
			});
		}

		getOrganizationList() {
			var getTemplate = this.getTemplate('organization-list');
			getTemplate.then(function() {
				$.ajax({
					method : 'get',
					url : '/api/1/organization',
					success : function (data, status, jqXHR) {
						var source = $('#template-organization-list').html();
						var template = Handlebars.compile(source);
						var html = template({list:data});
						$('#organization-list').append(html);
					}
				});
			});
		}
	}

	export class Main {
		init() {
			$('.button-collapse').sideNav();

			Handlebars.registerHelper('formatDatetime', function(time : string) {
				var startAt = moment(time);
				var zone = moment().zone();
				return startAt.zone(zone).format('YYYY/M/D H:mm')+'〜';
			});

			Handlebars.registerHelper('getColorClassByDate', function(time : string) {
				var startAt = moment(time);
				if ( startAt.isSame(moment(), 'day') ) {
					return 'amber lighten-5';
				}
				if ( startAt.isBefore(moment(), 'day') ) {
					return 'grey lighten-3';
				}
				return '';
			});

			var api = new Api();
			api.getEventList();
			api.getOrganizationList();
		}
	}

}

$(function() {
	var main = new GCPUG.Main();
	main.init();
});