$(document).ready(function(){
	var nocm = $("#nocm").val()
	
	$("#nocm").focus(function(){
		var value = $("#nocm").val();
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
			//nocm = "";
		} else {
			$("#datapasien").html("No. CM tidak lengkap");
			//nocm = "";
		} 
	});
	
	$("#form1").on("keyup", "#nocm", function(){
		
		var value = $("#nocm").val();
		
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
		} else if (value.length < 8){
			$("#datapasien").html("No. CM tidak lengkap");
		} else {
			$("#nocm").prop("disabled", true);
			nocm = value;
			$.get("/getcm", {nocm: value}, function(data){
				$("#datapasien").html(data);
				$("#nocm").prop("disabled", false);
				
				
			})
		}
	});
	
	$("#infouser").html(function(){
		$.get("/getinfo", function(data){
			$("#infouser").html(data)
		})
		
		
	});
	
	
})