'use strict';

// Declare app level module which depends on filters, and services
angular.module('bitnel', [
  'ui.router',
  'bitnel.controllers'
]).
run(
  [ '$rootScope', '$state', '$stateParams',
    function ($rootScope, $state, $stateParams) {

    // It's very handy to add references to $state and $stateParams to the $rootScope
    // so that you can access them from any scope within your applications.For example,
    // <li ng-class="{ active: $state.includes('contacts.list') }"> will set the <li>
    // to active whenever 'contacts.list' or one of its decendents is active.
    $rootScope.$state = $state;
    $rootScope.$stateParams = $stateParams;
    }
  ]
).
config(['$urlRouterProvider', '$stateProvider', '$locationProvider', function($urlRouterProvider, $stateProvider, $locationProvider) {
  $locationProvider.html5Mode(true);

  //$urlRouterProvider.otherwise('/index');

  $stateProvider
    .state('orders', {
      url: '/orders',
      abstract: true,
      templateUrl: '/views/orders.html',
    })
    .state('orders.list', {
      url: '/all',
      templateUrl: '/views/orders.list.html',
      controller: 'OrdersListCtrl'
    })
    .state('orders.new', {
      url: '/new',
      templateUrl: '/views/orders.new.html',
      controller: 'OrdersNewCtrl'
    })
    .state('orders.detail', {
      url: '/:orderUuid',
      templateUrl: '/views/orders.detail.html',
      controller: 'OrdersDetailCtrl'
    })
    .state('markets', {
      url: '/markets',
      abstract: true,
      templateUrl: '/views/markets.html',
    })
    .state('markets.detail', {
      title: 'Lol',
      url: '/:currencyPair',
      templateUrl: '/views/markets.detail.html',
      controller: 'MarketsDetailCtrl'
    })
    .state('withdrawals', {
      url: '/withdrawals',
      abstract: true,
      templateUrl: '/views/withdrawals.html',
    })
    .state('withdrawals.list', {
      url: '/',
      templateUrl: '/views/withdrawals.list.html',
      controller: 'WithdrawalsListCtrl'
    });
}]);