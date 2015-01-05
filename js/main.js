var webChat = angular.module('webChat', ['ngRoute']);

webChat.config(['$routeProvider',
  function ($routeProvider) {
      $routeProvider.
        when('/', {
            templateUrl: 'views/chatView.html',
            controller: 'ChatCtrl'
        }).
        otherwise({
            redirectTo: '/'
        });
  }]);


webChat.controller('ChatCtrl', function ($scope, $http) {
    $scope.posts = "";
    $scope.updatePosts = function () {
        $http.get('http://127.0.0.1:8080/chat/hello').
      success(function (data, status, headers, config) {
          $scope.posts = data;
      }).
      error(function (data, status, headers, config) {
          // called asynchronously if an error occurs
          // or server returns response with an error status.
      })
    };


    $scope.sendPost = function () {
        $http.get('http://127.0.0.1:8080/chat/hello', {msg:$scope.sendText}).
         success(function (data, status, headers, config) {
             $scope.updatePosts();
             $scope.sendText="";
         }).
         error(function (data, status, headers, config) {
             // called asynchronously if an error occurs
             // or server returns response with an error status.
         })
    }

    setInterval($scope.updatePosts, 1000);
    
});