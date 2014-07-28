'use strict';

angular.module('bitnel.controllers', [])
  .controller('OrdersListCtrl', ['$scope', '$http', '$route', function($scope, $http, $route) {
    $scope.title = "Orders";
    $http.get('api/orders').success(function(data) {
      //$scope.orders = data;
    });

    $scope.orders = [];

    for (var i = 0; i < 100; i++) {
      $scope.orders.push({
        uuid: guid(),
        quantity: 123+i,
        initial_quantity: 123,
        price: 299+i,
        status: "OPEN"
      })
    }

    $scope.$route = $route;
  }])
  .controller('OrdersNewCtrl', ['$scope', function($scope) {

  }])
  .controller('OrdersDetailCtrl', ['$scope', '$http', '$route', '$stateParams', function($scope, $http, $route, $stateParams) {
    $scope.orderId = $stateParams.orderUuid;

    $http.get('fd/api/orders/'+$stateParams.orderUuid).success(function(data) {
      //$scope.order = data;
    });

    $scope.order = {
      uuid: "lol",
      quantity: 123,
      initial_quantity: 123,
      price: 299,
      status: "OPEN"
    };

    $scope.$route = $route;
  }])
  .controller('MarketsDetailCtrl', ['$scope', '$http', '$route', '$stateParams', function($scope, $http, $route, $stateParams) {
    $scope.currencyPair = $stateParams.currencyPair;

    $http.get('fd/api/orders/'+$stateParams.orderUuid).success(function(data) {
      //$scope.order = data;
    });

    $scope.order = {
      uuid: "lol",
      quantity: 123,
      initial_quantity: 123,
      price: 299,
      status: "OPEN"
    };

    var canvas = document.getElementById('chart');
    canvas.width = canvas.offsetWidth;
    canvas.height = 300;
    var ctx = canvas.getContext('2d');

    canvas.addEventListener('mousemove', function(evt) {
      var mousePos = getMousePos(canvas, evt);
      ctx.clearRect(0, 0, window.innerWidth-200, 300);

      ctx.beginPath();
      ctx.moveTo(mousePos.x, 0);
      ctx.lineTo(mousePos.x, 300);
      ctx.stroke();
      console.log(mousePos.x);


      ctx.font = "100px Arial";
      ctx.fillStyle = 'rgba(0, 0, 0, .1)';
      ctx.fillText($scope.currencyPair, 150, 130);


      ctx.font = "50px Arial";
      ctx.fillText("© 2014 Bitnel", 150, 180);
    }, false);

    

  function getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect(), root = document.documentElement;

    // return relative mouse position
    var mouseX = evt.clientX - rect.left - root.scrollLeft;
    var mouseY = evt.clientY - rect.top - root.scrollTop;
    return {
      x: mouseX,
      y: mouseY
    };
  }


    $scope.$route = $route;
  }])
  .controller('WithdrawalsNewCtrl', ['$scope', function($scope) {

  }])
  .controller('WithdrawalsListCtrl', ['$scope', function($scope) {

  }]);

var guid = (function() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
               .toString(16)
               .substring(1);
  }
  return function() {
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
           s4() + '-' + s4() + s4() + s4();
  };
})();