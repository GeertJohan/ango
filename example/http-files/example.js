var app = angular.module('chatExample', ['chatService']);

app.config(function(chatServiceProvider) {
	chatServiceProvider.setUri("");
	chatServiceProvider.handle('newMessage', function(message) {
		console.log('Have a new message');
	});
});

app.controller('mainCtrl', function($scope, chatService) {
	$scope.foo = "bar";
	console.log(chatService.getWsPath());
})