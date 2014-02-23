var chat = angular.module('chatExample', ['ango-chatService']);

chat.config(function(chatServiceProvider) {
	chatServiceProvider.setDebug(true);
});

chat.controller('mainCtrl', function($scope, chatService) {
	$scope.foo = "bar";
	console.log(chatService.getWsPath());
});