var chat = angular.module('chatExample', ['ango-chatService']);

chat.config(function(chatServiceProvider) {
	chatServiceProvider.setUri("");
	chatServiceProvider.handle('newMessage', function(message) {
		console.log('Have a new message');
	});
});

chat.controller('mainCtrl', function($scope, chatService) {
	$scope.foo = "bar";
	console.log(chatService.getWsPath());
});